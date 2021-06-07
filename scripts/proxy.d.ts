//
// See https://www.typescriptlang.org/docs/handbook/release-notes/typescript-2-9.html#import-types
//

type Middleware = import("./middleware/types").Middleware;

declare namespace proxy {
  function onInit(initFn: () => void);

  function onDestroy(initFn: () => void);

  function addReceivedMessageMiddleware(...fn: Middleware[]);

  function addSentMessageMiddleware(...fn: Middleware[]);

  function addReceiveMessageMiddleware(...fn: Middleware[]);

  function addSendMessageMiddleware(...fn: Middleware[]);
}
