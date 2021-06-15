import {Message, Middleware} from "./types";

export function getLogger(prefix: string): Middleware {
  return function (message: string | Message): string | Message {
    console.log(prefix, message);
    return message;
  }
}
