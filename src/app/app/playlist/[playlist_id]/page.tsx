import { SongsTable } from "@/components/SongsTable";
import { PlaylistRecord, SongRecord } from "@/interfaces";
export default async function PlayistPage({
  params,
}: {
  params: { playlist_id: string };
}) {
  const playlistDataReq = fetch(
    `${process.env["BACKEND_URL"]}/playlist/${params.playlist_id}`,
    { cache: "no-store" }
  );
  const tracksReq = fetch(
    `${process.env["BACKEND_URL"]}/playlist/${params.playlist_id}/songs`,
    { cache: "no-store" }
  );

  const [playlistDataRes, tracksRes] = await Promise.all([
    playlistDataReq,
    tracksReq,
  ]);
  const playlistData: PlaylistRecord = await playlistDataRes.json();
  const tracks: SongRecord[] = (await tracksRes.json())["songs"];

  return (
    <div className="h-screen max-h-full overflow-scroll">
      <h1>{playlistData.title}</h1>
      <SongsTable songs={tracks} />
    </div>
  );
}
