import { NextRequest, NextResponse } from "next/server";

const SESSION_MAX_AGE = 7 * 24 * 60 * 60;

export async function GET(req: NextRequest) {
  const token = req.nextUrl.searchParams.get("token");
  if (!token) {
    return NextResponse.redirect(new URL("/login", req.url));
  }

  const res = NextResponse.redirect(new URL("/home", req.url));
  res.cookies.set("session", token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
    maxAge: SESSION_MAX_AGE,
  });
  return res;
}
