from django.apps import AppConfig

from django.http import HttpResponse 

def index(request):
    return HttpResponse("Hello, world. You're at the polls index.")

class PollsConfig(AppConfig):
    name = 'polls'
