import {getLogger} from "./middleware/log";

proxy.addReceivedMessageMiddleware(getLogger("<--"));

proxy.addSentMessageMiddleware(getLogger("-->"));

let timeoutID = 0;
proxy.onInit(function () {
  timeoutID = setTimeout(() => {
    console.log("5s elapsed.");
  }, 5000);
});

proxy.onDestroy(function () {
  clearTimeout(timeoutID);
});
