from __future__ import unicode_literals

from django.db import models


class LentaNews(models.Model):
    link = models.URLField(primary_key=True)
    title = models.CharField(max_length=255)
    description = models.TextField()
    category = models.CharField(max_length=50, db_index=True)
    pub_date = models.DateTimeField(db_index=True)

    class Meta:
        ordering = ['category', '-pub_date']


