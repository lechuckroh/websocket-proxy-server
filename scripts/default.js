import {logMiddleware} from "./middleware/log";

proxy.addReceivedMessageMiddleware(logMiddleware("<--"));

proxy.addSentMessageMiddleware(logMiddleware("-->"));

proxy.onInit(function() {
  console.log("onInit");
});

proxy.onDestroy(function() {
  console.log("onDestroy");
});
