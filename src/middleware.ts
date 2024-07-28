import { NextResponse, type NextRequest } from "next/server";

export async function middleware(request: NextRequest) {
  // Reverse proxy for backend
  if (request.nextUrl.pathname.startsWith("/oauth")) {
    return NextResponse.rewrite(
      new URL(`${process.env["BACKEND_URL"]}/login`, request.url)
    );
  }

  // Check if the user is authenticated
  // only run checks for login and app pages
  if (
    !request.nextUrl.pathname.startsWith("/login") &&
    !request.nextUrl.pathname.startsWith("/app")
  ) {
    return;
  }

  let valid = false;
  if (request.cookies.get("pb_auth")?.value) {
    try {
      const parsed = JSON.parse(
        decodeURIComponent(request.cookies.get("pb_auth")?.value ?? "")
      );
      const validationResponse = await fetch(
        `${process.env["BACKEND_URL"]}/validate?user_id=${parsed?.user_id}&token=${parsed?.access_token}`
      );
      valid = (await validationResponse.json())["valid"];
    } catch (e) {
      console.error(e);
    }
  }

  if (valid && request.nextUrl.pathname.startsWith("/login")) {
    return Response.redirect(new URL("/app", request.url));
  }
  if (!valid && request.nextUrl.pathname.startsWith("/app")) {
    return Response.redirect(new URL("/login", request.url));
  }
}

export const config = {
  matcher: ["/app/:path*", "/login", "/oauth"],
};
