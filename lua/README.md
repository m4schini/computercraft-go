Client
---

For this library to work we need some client side code on the computercraft devices.

One Call should look like the following:
1. Receive websocket message with a lua expression
2. execute lua expression
3. send websocket message with wrapped (`{}`) returns of the lua expression