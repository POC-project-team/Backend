import requests as r


def random_string():
    import string
    import random
    return ''.join(random.choice(string.ascii_letters) for i in range(10))

url = "http://localhost:60494/"

def main():
    print(r.get(url=f"{url}users").json())
    for _ in range (10):
        login = random_string()
        password = login
        r.post(url=f'{url}signup', json={'login': login, 'password': password})
        token = r.post(url=f'{url}auth', json={'login': login, 'password': password})
        print(token.json()['token'])

        tag_name = random_string()

        create_link = f"{url}{token.json()['token']}/{tag_name}/tag"

        create = r.post(url=create_link, json={'tagName': "test"})
        print(create.json())

        add_note = f"{url}{token.json()['token']}/{tag_name}/note"
        r.post(url=add_note, json={'note': "test"})
        r.post(url=add_note, json={'note': "test123"})

        # get notes
        print(r.get(url=f"{url}{token.json()['token']}/{tag_name}/notes").json())
        


if __name__ == '__main__':
    main()
