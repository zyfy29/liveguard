import requests

from front.config import settings

# DYNACONF_API_HOST, DYNACONF_API_HOST
base_url = f'http://{settings.api_host}:{settings.api_port}'


def get_lives(owner_id, next_time=0):
    params = {
        'owner_id': owner_id
    }
    if next_time:
        params['next_time'] = next_time
    res = requests.get(f'{base_url}/pocket/live', params=params)
    data = res.json()
    return data['data']['lives'], data['data']['next_time']


def create_task(live_id):
    res = requests.post(f'{base_url}/task', json={'live_id': live_id})
    return res.json()


def retry_task(task_id):
    res = requests.post(f'{base_url}/task/retry', json={'task_id': task_id})
    return res.json()


def delete_task(task_id):
    res = requests.delete(f'{base_url}/task/{task_id}')
    return res.json()


def get_members():
    res = requests.get(f'{base_url}/pocket/member')
    return res.json()['data']


def get_tasks(status=''):
    res = requests.get(f'{base_url}/task', params={'status': status})
    return res.json()


def get_task_detail(task_id):
    res = requests.get(f'{base_url}/task/{task_id}')
    return res.json()


def get_my_name():
    res = requests.get(f'{base_url}/pocket/me')
    return res.json()
