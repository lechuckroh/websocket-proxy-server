export type Message = { [key: string]: any };

export type Middleware = (
  message: string | Message,
  ...rest: unknown[]
) => string | Message | null;
