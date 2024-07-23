"use client";
import { useAuth } from "@/components/PocketBaseAuthProvider";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import Link from "next/link";
import { useRouter } from "next/navigation";
import PocketBase, {
  ClientResponseError,
  RecordAuthResponse,
  RecordModel,
} from "pocketbase";
import { useEffect, useState } from "react";

export default function LoginForm() {
  const pb = new PocketBase(process.env.NEXT_PUBLIC_POCKETBASE_URL);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();
  const [authData, setAuthData] =
    useState<null | RecordAuthResponse<RecordModel>>(null);
  const { isAuthenticated } = useAuth();

  if (isAuthenticated) router.push("/app");

  const spotifyOAuth = async (event: { preventDefault: () => void }) => {
    event.preventDefault();
    try {
      const authData = await pb.collection("users").authWithOAuth2({
        provider: "spotify",
        scopes: ["user-read-email", "playlist-modify-private"],
      });
      console.log("Logged in successfully:", authData);
      console.log(authData?.meta?.accessToken);
      console.log(authData?.meta?.refreshToken);
      console.log(authData.record.id);

      console.log(authData?.meta?.email || authData?.meta?.rawUser?.email);

      const data = {
        user_id: authData.record.id,
        spotify_id: authData?.meta?.id,
        spotify_username: authData?.meta?.username,
        spotify_email: authData?.meta?.email || authData?.meta?.rawUser?.email,
        accessToken: authData?.meta?.accessToken,
        refreshToken: authData?.meta?.refreshToken,
        expiry: authData?.meta?.expiry,
        display_name: authData?.meta?.rawUser?.display_name,
        avatar_url: authData?.meta?.avatarUrl,
      };

      await pb.collection("users").update(authData.record.id, data);

      router.push("/app");
    } catch (err) {
      console.error("Failed to log in:", err);
      setError("Login failed");
    }
  };
  return (
    <section className="flex flex-col items-center justify-center min-h-screen">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-2xl">Login</CardTitle>
          <CardDescription>
            Link your Spotify account to get started
          </CardDescription>
        </CardHeader>
        <form onSubmit={spotifyOAuth}>
          <CardContent className="grid gap-4">
            {error && <p className="text-red-500">{error}</p>}
          </CardContent>
          <CardFooter className="flex flex-col items-center gap-1">
            <Button type="submit" className="w-full">
              Sign in with Spotify
            </Button>
          </CardFooter>
        </form>
      </Card>
    </section>
  );
}
