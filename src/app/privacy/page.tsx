import Link from "next/link";

export default function Privacy() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <h1 className="text-6xl font-semibold">Privacy Policy</h1>
      <div className=" text-center">
        Daylist Rewind is an open souce app powered by the Spotify Web API. By
        choosing to use this app, you agree share the following information with
        us. Here's what you need to know about what we collect and how we use
        it:<br></br>
        <h2 className="text-2xl font-semibold">
          The permissions you are granting us access to
        </h2>
        <ul>
          <li className="">
            <span className="font-bold">
              The email associated with your Spotify account
            </span>
            . We only use this for identification and we will never send you
            emails
          </li>
          <li>
            <span className="font-bold">
              Your Spotify username / public profile
            </span>
          </li>
          <li>
            The ability to <span className="font-bold">write</span> to but not
            read from your <span className="font-bold">private playlists</span>.
            We need this to be able to export a playlist to your account.
          </li>
        </ul>
        <h2>The information we store</h2>
        <ul>
          <li>
            <span className="font-bold">Your Spotify profile</span>
          </li>
          <li>
            <span className="font-bold">
              The daylists that Spotify creates for you
            </span>
          </li>
        </ul>
        <div>
          You may at any time request the deletion of your data in your profile
          settings. This will delete our access keys from our database as well
          as any data we have stored about you. This effectively prevents us
          from accessing your Spotify account. However, you may also revoke our
          permissions entirely by following{" "}
          <Link
            className="underline"
            href="https://support.spotify.com/us/article/spotify-on-other-apps/"
          >
            this guide
          </Link>{" "}
          on how to do so.
        </div>
      </div>
    </main>
  );
}
