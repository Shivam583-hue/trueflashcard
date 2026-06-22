import { Code, ConnectError } from "@connectrpc/connect";

export function isUnauthenticated(err: unknown): boolean {
  return ConnectError.from(err).code === Code.Unauthenticated;
}
