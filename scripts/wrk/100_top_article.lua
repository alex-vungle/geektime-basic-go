wrk.method="GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["User-Agent"] = "PostmanRuntime/7.32.3"
-- 记得修改这个，你在登录页面登录一下，然后复制一个过来这里
wrk.headers["Authorization"]="Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjIwMDExNjMsIlVpZCI6MiwiU3NpZCI6ImZkY2Q4NDRlLWM5N2YtNDM0NC1hYzZlLWZmOTQ1OWNhMzFlMiIsIlVzZXJBZ2VudCI6IlBvc3RtYW5SdW50aW1lLzcuMzkuMCJ9.XJYrQgZC9X3TLUEc8pdGkE_YtgV5OBrLbHtKw3xOMVbbKfGQ4Js7zCu_hg4tYmMHr1kLk3esN7x5BcvwAIw1Zg"