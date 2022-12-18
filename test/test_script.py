import requests as r


def random_string():
    import string
    import random
    return ''.join(random.choice(string.ascii_letters) for i in range(10))


def main():
    url = "http://localhost:60494/"

    print(r.get(url=f"{url}users").json())
    for _ in range(1):
        login = random_string()
        password = login
        res = r.post(url=f'{url}signup', json={'login': login, 'password': password})
        print(res)
        print(login, password)
        token = r.post(url=f'{url}auth', json={'login': login, 'password': password})
        print(token)
        print(f"token: {token.json()['token']}")
        tag_name = random_string()

        create_link = f"{url}{token.json()['token']}/{tag_name}/tag"
        create = r.post(url=create_link, json={'tagName': "test"})
        print(f"Just created tag: {create.json()}")

        # get all tags from the user
        create_link = f"{url}{token.json()['token']}/tags"
        tags = r.get(url=create_link).json()
        print(f"all tags from the user: {tags}")

        add_note = f"{url}{token.json()['token']}/{tag_name}/note"
        r.post(url=add_note, json={'note': "test"})
        r.post(url=add_note, json={'note': "test123"})
        r.post(url=add_note, json={'note': login})

        # # get notes
        notes = r.get(url=f"{url}{token.json()['token']}/{tag_name}/notes").json()
        print(f"All notes from the user {notes}")

    print("Use-case has been successfully completed")
        

if __name__ == '__main__':
    main()
