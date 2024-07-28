"use client";

import { PlaylistRecord } from "@/interfaces";
import { formatDate, removeDaylist } from "@/util";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import DOMPurify from "dompurify";
import { useRouter } from "next/navigation";

function convertToSpotifyLinks(description: string) {
  const linkRegex = /<a href="([^"]+)">([^<]+)<\/a>/g;
  let result = description;
  let matches;

  while ((matches = linkRegex.exec(description)) !== null) {
    const link = matches[1];
    const text = matches[2];

    if (link.startsWith("spotify:")) {
      const spotifyLink = `https://open.spotify.com/${link
        .replace(/:/g, "/")
        .replace("spotify/", "")}`;
      result = result.replace(
        matches[0],
        `<a class="underline" href="${spotifyLink}">${text}</a>`
      );
    }
  }

  return result;
}

function sanitize(html: string) {
  // dumb as fuck nextjs workarounds
  if (typeof window == "undefined") {
    return html;
  }
  const DOMPurifyServer = DOMPurify(window)
  return DOMPurifyServer.sanitize(
    convertToSpotifyLinks(html)
  )
}


export function Playlists({ playlists }: { playlists: PlaylistRecord[] }) {
  const router = useRouter();

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-36">Date</TableHead>
          <TableHead>Title</TableHead>
          <TableHead>Description</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {playlists.map((playlist: PlaylistRecord) => (
          <TableRow
            key={playlist.id}
            onClick={() => router.push(`/app/playlist/${playlist.id}`)}
          >
            <TableCell>{formatDate(playlist.created)}</TableCell>
            <TableCell>{removeDaylist(playlist.title)}</TableCell>
            <TableCell
              dangerouslySetInnerHTML={{
                __html: sanitize(
                  convertToSpotifyLinks(playlist.description)
                ),
              }}
            ></TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
