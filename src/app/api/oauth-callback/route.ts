import { cookies, headers } from "next/headers";
import { redirect } from "next/navigation";

// forwards the oauth callback to the server at localhost:8080/callback
export async function GET(request: Request) {
  const url = new URL(request.url);
  const code = url.searchParams.get("code");
  const state = url.searchParams.get("state");
  // const ip = headers().get('X-Forwarded-For')

  const response = await fetch(
    `${process.env["BACKEND_URL"]}/callback?code=${code}&state=${state}`,
  );

  const data = await response.json();

  // stores the access token in a cookie
  const accessToken = data.access_token;
  const expires = new Date(Date.now() + data.expires_in * 1000);

  const pbAuthCookie = encodeURIComponent(JSON.stringify(data));
  cookies().set("pb_auth", pbAuthCookie, {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    maxAge: 60 * 60, // One hour
    path: "/",
  });
  redirect("/login");

  // return Response.json(data);
}
