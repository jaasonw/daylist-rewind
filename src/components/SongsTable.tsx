"use client";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { SongRecord } from "@/interfaces";
import { msToMinSec } from "@/util";

export function SongsTable({ songs }: { songs: SongRecord[] }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-3/12">Title</TableHead>
          <TableHead className="w-2/12">Artist</TableHead>
          <TableHead>Album</TableHead>
          <TableHead>Duration</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {songs.map((track: SongRecord) => (
          <TableRow
            key={track.song_id}
            onClick={() => {
              window.location.href = `https://open.spotify.com/track/${track.song_id}`;
            }}
            className="cursor-pointer"
          >
            <TableCell>{track.name}</TableCell>
            <TableCell>{track.artist}</TableCell>
            <TableCell>{track.album}</TableCell>
            <TableCell>{msToMinSec(track.duration)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
