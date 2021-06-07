import {Message, Middleware} from "./types";

export function logMiddleware(prefix: string): Middleware {
  return function (message: string | Message): string | Message {
    console.log(prefix, message);
    return message;
  }
}
