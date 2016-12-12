from django.conf.urls import url
from digest import views


urlpatterns = [
    url(r'^send_digest_to_email$', views.send_digest_to_email),
    url(r'^check_task_status$', views.check_task_status),
    url(r'^$', views.index),
]
