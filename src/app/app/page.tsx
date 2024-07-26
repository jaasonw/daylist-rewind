import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { PlaylistRecord, UserRecord } from "@/interfaces";
import DOMPurify from "dompurify";
import { JSDOM } from "jsdom";
import { cookies } from "next/headers";
import PocketBase from "pocketbase";

const window = new JSDOM("").window;
const DOMPurifyServer = DOMPurify(window);

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

function formatDate(dateString: string) {
  const date = new Date(dateString);
  const mm = String(date.getMonth() + 1).padStart(2, "0"); // Months are zero-based, so we add 1
  const dd = String(date.getDate()).padStart(2, "0");
  const yyyy = date.getFullYear();

  return `${mm}-${dd}-${yyyy}`;
}

function removeDaylist(title: string) {
  return title.replaceAll("daylist • ", "");
}

export default async function Dashboard() {
  const pb = new PocketBase(process.env.POCKETBASE_URL);
  await pb.admins.authWithPassword(
    process.env["ADMIN_USER"] ?? "",
    process.env["ADMIN_PASSWORD"] ?? ""
  );

  const cookie = cookies()?.get("pb_auth");
  const cookieData = JSON.parse(decodeURIComponent(cookie?.value ?? ""));
  // console.log(cookieData);
  const userId = cookieData.user_id;

  const userData: UserRecord = await pb
    .collection("users")
    .getFirstListItem(`username="${userId}"`);
  // console.log(userData);
  // console.log(cookieData);
  // const userId = cookieData;

  const playlistsResponse = await fetch(
    `${process.env["BACKEND_URL"]}/user/playlists/${userData.id}`
  );
  const playlists: PlaylistRecord[] = await playlistsResponse.json();
  console.log(playlists);
  return (
    <>
      <div className="flex items-center">
        <h1 className="text-lg font-semibold md:text-2xl">Welcome</h1>
      </div>
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
            <TableRow key={playlist.id}>
              <TableCell>{formatDate(playlist.created)}</TableCell>
              <TableCell>{removeDaylist(playlist.title)}</TableCell>
              <TableCell
                dangerouslySetInnerHTML={{
                  __html: DOMPurifyServer.sanitize(
                    convertToSpotifyLinks(playlist.description)
                  ),
                }}
              ></TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      {/* <div
        className="flex flex-1 items-center justify-center rounded-lg border border-dashed shadow-sm"
        x-chunk="dashboard-02-chunk-1"
      >
      </div> */}
    </>
  );
}
