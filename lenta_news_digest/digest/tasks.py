# coding: utf-8
import os
import urllib2
from lxml import etree
import logging
from dateutil.parser import parse as date_parse
from datetime import datetime, timedelta
from StringIO import StringIO

from django.conf import settings
from django.core.mail import EmailMessage
import xhtml2pdf.pisa as pisa
from celery.task import periodic_task, task

from digest.models import LentaNews


log = logging.getLogger('digest')


@periodic_task(run_every=timedelta(minutes=5))
def fetch_lenta_rss():
    rss_url = 'https://lenta.ru/rss/news'
    page = urllib2.urlopen(rss_url)
    rss = etree.parse(page).getroot()
    items = rss.iterfind('.//item')
    fresh_news_count = 0
    fresh_news = []
    for item in items:
        link = item.find('link').text
        try:
            news_entry = LentaNews.objects.get(link=link)
        except LentaNews.DoesNotExist:
            try:
                title = item.find('title').text
                description = item.find('description').text
                category = item.find('category').text
                pub_date = date_parse(item.find('pubDate').text)
                fresh_news.append(LentaNews(link=link,
                                            title=title,
                                            category=category,
                                            description=description,
                                            pub_date=pub_date))
                fresh_news_count += 1
            except Exception as e:
                log.error('fetch_lenta_rss: %s: %s' % (link, e))
    LentaNews.objects.bulk_create(fresh_news)
    if fresh_news_count > 0:
        log.info('fetch_lenta_rss: %d fresh news' % fresh_news_count)


@task
def gen_pdf_and_send_mail(html, email, date_from_string, date_to_string):
    result = StringIO()
    path = os.path.join(settings.STATICFILES_DIRS[0], 'fonts', 'whatever')  # hack for unicode
    pdf = pisa.pisaDocument(StringIO(html.encode('utf-8')), result, path=path)

    # Send mail
    if date_from_string == date_to_string:
        period_text = u'За %s' % date_from_string
    else:
        period_text = u'За период с %s по %s' % (date_from_string, date_to_string)
    subject = u'Дайджест новостей %s' % period_text
    body = 'Дайджест в прикреплённом pdf-документе'
    mail = EmailMessage(
        subject=subject,
        body=body,
        from_email='Django application',
        to=[email],
    )
    mail.attach(filename='digest.pdf', content=result.getvalue(), mimetype='application/pdf')
    try:
        result = mail.send(fail_silently=False)
        if result == 1:
            response = {'success':True, 'msg':u'Письмо отправлено успешно'}
            return response
        else:
            response = {'success':False,
                        'msg':u'Что-то пошло не так: результат mail.send = %s' % result} 
            return response
    except Exception, e:
        response = {'success':False, 'msg':u'Произошла ошибка: %s' % e} 
        return response



