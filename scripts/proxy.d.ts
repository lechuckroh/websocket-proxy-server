//
// See https://www.typescriptlang.org/docs/handbook/release-notes/typescript-2-9.html#import-types
//

import {ResponseFunc} from "./middleware/types";

type Middleware = import("./middleware/types").Middleware;

declare namespace proxy {
  function onInit(initFn: (responseToBackend: ResponseFunc, responseToClient: ResponseFunc) => void);

  function onDestroy(destroyFn: () => void);

  function addReceivedMessageMiddleware(...fn: Middleware[]);

  function addSentMessageMiddleware(...fn: Middleware[]);

  function addResponseToBackendMessageMiddleware(...fn: Middleware[]);

  function addResponseToClientMessageMiddleware(...fn: Middleware[]);
}
