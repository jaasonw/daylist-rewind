import DOMPurify from "dompurify";

export function formatDate(dateString: string) {
  const date = new Date(dateString);
  const mm = String(date.getMonth() + 1).padStart(2, "0"); // Months are zero-based, so we add 1
  const dd = String(date.getDate()).padStart(2, "0");
  const yyyy = date.getFullYear();

  return `${mm}-${dd}-${yyyy}`;
}

export function removeDaylist(title: string) {
  return title.replaceAll("daylist • ", "");
}

export function msToMinSec(ms: number) {
  const totalSeconds = Math.floor(ms / 1000);
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;

  // Format minutes and seconds to be two digits
  const formattedMinutes = minutes.toString().padStart(2, "0");
  const formattedSeconds = seconds.toString().padStart(2, "0");

  return `${formattedMinutes}:${formattedSeconds}`;
}

export function convertToSpotifyLinks(description: string) {
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
        `<a class="underline" href="${spotifyLink}">${text}</a>`,
      );
    }
  }

  return result;
}

export function sanitize(html: string) {
  // dumb as fuck nextjs workarounds
  if (typeof window == "undefined") {
    return html;
  }
  const DOMPurifyServer = DOMPurify(window);
  return DOMPurifyServer.sanitize(convertToSpotifyLinks(html));
}
