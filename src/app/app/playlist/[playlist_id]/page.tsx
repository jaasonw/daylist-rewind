import { ExportPlaylistButton } from "@/components/ExportPlaylistButton";
import { SongsTable } from "@/components/SongsTable";
import { Button } from "@/components/ui/button";
import { PlaylistRecord, SongRecord } from "@/interfaces";
import { convertToSpotifyLinks, sanitize } from "@/util";
import { Plus } from "lucide-react";
export default async function PlayistPage({
  params,
}: {
  params: { playlist_id: string };
}) {
  const playlistDataReq = fetch(
    `${process.env["BACKEND_URL"]}/playlist/${params.playlist_id}`,
    { cache: "no-store" },
  );
  const tracksReq = fetch(
    `${process.env["BACKEND_URL"]}/playlist/${params.playlist_id}/songs`,
    { cache: "no-store" },
  );

  const [playlistDataRes, tracksRes] = await Promise.all([
    playlistDataReq,
    tracksReq,
  ]);
  const playlistData: PlaylistRecord = await playlistDataRes.json();
  const tracks: SongRecord[] = (await tracksRes.json())["songs"];

  return (
    <div className="h-screen max-h-full overflow-scroll p-4">
      <div className="flex flex-col justify-end h-96 mb-10">
        <div className="flex flex-col gap-3">
          <h1 className="font-bold text-6xl">{playlistData.title}</h1>
          <h2
            className="text-lg"
            dangerouslySetInnerHTML={{
              __html: sanitize(convertToSpotifyLinks(playlistData.description)),
            }}
          ></h2>
          <div>
            <ExportPlaylistButton />
          </div>
        </div>
      </div>
      <SongsTable songs={tracks} />
    </div>
  );
}
