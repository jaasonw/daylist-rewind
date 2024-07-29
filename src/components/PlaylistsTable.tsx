"use client";

import { PlaylistRecord } from "@/interfaces";
import { convertToSpotifyLinks, formatDate, removeDaylist, sanitize } from "@/util";
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
            className="cursor-pointer"
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
