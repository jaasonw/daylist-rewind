export const dynamic = "force-dynamic";
import { Playlists } from "@/components/PlaylistsTable";
import { PlaylistRecord, UserRecord } from "@/interfaces";
import { cookies } from "next/headers";

export default async function Dashboard() {
  const cookie = cookies()?.get("pb_auth");
  const cookieData = JSON.parse(decodeURIComponent(cookie?.value ?? ""));
  const userId = cookieData.user_id;

  const userData: UserRecord = await fetch(
    `${process.env["BACKEND_URL"]}/user/${userId}?access_token=${cookieData.access_token}`,
  ).then((res) => res.json());
  const playlistsResponse = await fetch(
    `${process.env["BACKEND_URL"]}/user/playlists/${userData.id}`,
  );
  const playlists: PlaylistRecord[] = await playlistsResponse.json();
  return (
    <>
      <div className="flex items-center">
        <h1 className="text-lg font-semibold md:text-2xl">Welcome</h1>
      </div>
      <Playlists playlists={playlists} />
      {/* <div
        className="flex flex-1 items-center justify-center rounded-lg border border-dashed shadow-sm"
        x-chunk="dashboard-02-chunk-1"
      >
      </div> */}
    </>
  );
}
