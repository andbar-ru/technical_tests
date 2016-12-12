# coding: utf-8
from datetime import date, time
from dateutil.parser import parse as date_parse
import json
import logging
import re
from collections import OrderedDict

from django.shortcuts import render
from django.template.loader import render_to_string
from django.http import HttpResponse
from celery.result import AsyncResult

from digest.models import LentaNews
from digest.tasks import gen_pdf_and_send_mail


log = logging.getLogger('digest')


def json_response(response):
    json_string = json.dumps(response, ensure_ascii=False).encode('utf-8')
    return HttpResponse(json_string)


def index(request):
    categories = (LentaNews.objects.order_by('category')
                                   .values_list('category', flat=True)
                                   .distinct())
    news_categories = OrderedDict()
    for category in categories:
        news = LentaNews.objects.filter(category=category)[:5]
        if len(news) > 0:
            news_categories[category] = news

    today = date.today().isoformat()
    min_date = (LentaNews.objects.order_by('pub_date')
                                 .values_list('pub_date', flat=True)[0]
                                 .date().isoformat())
    max_date = (LentaNews.objects.order_by('-pub_date')
                                 .values_list('pub_date', flat=True)[0]
                                 .date().isoformat())
    context = {
        'news_categories': news_categories,
        'categories': categories,
        'today': today,
        'min_date': min_date,
        'max_date': max_date,
    }
    return render(request, 'index.djhtml', context)


def check_task_status(request):
    task_id = request.GET.get('task_id')
    if task_id:
        async_result = AsyncResult(task_id)
        log.info('%s %s' % (async_result.ready(), async_result.status))
        if async_result.ready():
            if async_result.status == 'SUCCESS':
                return json_response(async_result.result)
            else:
                return json_response({'success':False, 'msg':async_result.result})
        else:
            return json_response({'success':False, 'msg':'pending'})
    else:
        return json_response({'success':False, 'msg':'Пустой task_id'})


def send_digest_to_email(request):
    if request.method == 'POST':
        data = request.POST

        # validation
        email = data['email']
        if not email or not re.match(r'^[\w\d._%+-]+@[\w\d.-]+\.[\w]{2,}$', email):
            return json_response({'success':False, 'msg':u'Пустой или неправильный email'})
        date_from = data['date_from']
        if not date_from:
            return json_response({'success':False, 'msg':u'Пустая дата "с"'})
        date_to = data['date_to']
        if not date_to:
            return json_response({'success':False, 'msg':u'Пустая дата "до"'})
        try:
            date_from = date_parse(date_from)
        except:
            return json_response({'success':False, 'msg':u'Неправильная дата "с"'})
        try:
            date_to = date_parse(date_to)
        except:
            return json_response({'success':False, 'msg':u'Неправильная дата "до"'})
        if date_from > date_to:
            return json_response({'success':False, 'msg':u'Дата "с" не может быть больше даты "до"'})

        date_from = date_from.combine(date_from, time.min)
        date_to = date_to.combine(date_to, time.max)
        categories = data.getlist('category')
        if 'all' in categories:
            categories = (LentaNews.objects.order_by('category')
                                           .values_list('category', flat=True)
                                           .distinct())
        news_categories = OrderedDict()
        for category in categories:
            q = LentaNews.objects.filter(category=category,
                                         pub_date__gte=date_from,
                                         pub_date__lte=date_to)
            if q.count() > 0:
                news_categories[category] = q
        date_from_string = date_from.strftime('%d.%m.%Y')
        date_to_string = date_to.strftime('%d.%m.%Y')
        context = {
            'news_categories': news_categories,
            'date_from': date_from_string,
            'date_to': date_to_string,
        }
        html = render_to_string('digest.djhtml', context=context)
        task = gen_pdf_and_send_mail.apply_async((html, email, date_from_string, date_to_string))

        return json_response({'success':True, 'msg':task.task_id})
