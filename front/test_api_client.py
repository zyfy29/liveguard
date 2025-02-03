from api_client import get_lives


def test_get_data():
    res = get_lives(owner_id=63566)
    print(res)
    assert len(res.get('lives')) > 0
