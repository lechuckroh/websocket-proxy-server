export type Message = { [key: string]: any };

export type ResponseFunc = (message: string | Message) => void;

export type Middleware = (
  message: string | Message,
  response: ResponseFunc,
  ...rest: unknown[]
) => string | Message | null;
