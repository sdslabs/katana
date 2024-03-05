import requests


# if anything goes wrong then returns 0 at that instant otherwise returns 1 at last
def check_challenge(url):
    result = {}

    register_response = requests.post(
        url + "register", json={"username": "something", "voucher": ""}
    )

    # Check if the request was successful (status code 200)
    if register_response.status_code == 200:
        # Extract the value of the 'token' cookie
        token_cookie = register_response.cookies.get("token")

        # Print or use the token as needed
        result["Token Cookie"] = token_cookie
    else:
        result["Request failed with status code"] = register_response.status_code
        result["Response body"] = register_response.text
        return 0

    article_headers = {
        "Accept": "*/*",
        "Connection": "keep-alive",
        "Content-Type": "application/json",
        "Cookie": f"token={token_cookie}",
    }

    article_data = {"issue": "5"}

    article_response = requests.post(
        url + "article", json=article_data, headers=article_headers
    )

    if article_response.status_code == 200:
        result["Article response"] = article_response.text
        # print("Article response:", article_response.text)
    else:
        return 0

    # exploit ------ not needed -------

    # exploit = {"issue": "9a"}

    # exploit_response = requests.post(
    #     url + "article", json=exploit, headers=article_headers
    # )
    # if exploit_response.status_code == 200:
    #     result["Exploit response"] = exploit_response.text
    #     result["Flag"] = exploit_response.json()["content"]

    return 1
