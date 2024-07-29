"use server";

import { cookies } from "next/headers";

export async function GET(
  request: Request,
  { params }: { params: { playlist_id: string } }
) {
  const playlist_id = params.playlist_id
  const cookie = cookies()?.get("pb_auth");
  const cookieData = JSON.parse(decodeURIComponent(cookie?.value ?? ""));
  const userId = cookieData.user_id;
  const response = await fetch(`${process.env["BACKEND_URL"]}/playlist/${playlist_id}/export?username=${userId}&access_token=${cookieData.access_token}`)
  const data = await response.json()
  return new Response(`Hello ${JSON.stringify(data)}!`)
}