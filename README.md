Swagger Url: http://localhost:1323/swagger/index.html

For Simple Load Test Suit:

npm install -g load-tester

load-tester 5000

http://localhost:5000

Example Test Setup:

````
{
    "baseUrl": "http://localhost:1323",
    "duration": 15000,
    "connections": 100,
    "sequence": [
    { "method": "POST", "path": "/coupon/api/v1/give-transaction/d7513edf-511c-4515-a7e7-2920d50f1237"}
    ]
}
````
