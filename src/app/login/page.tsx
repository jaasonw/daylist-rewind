import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import Link from "next/link";
export default function LoginForm() {
  return (
    <section className="flex flex-col items-center justify-center min-h-screen">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-2xl">Login</CardTitle>
          <CardDescription>
            Link your Spotify account to get started
          </CardDescription>
        </CardHeader>
        <form>
          <CardContent className="grid gap-4">
            {/* {error && <p className="text-red-500">{error}</p>} */}
          </CardContent>
          <CardFooter className="flex flex-col items-center gap-1">
            <Link href="./oauth">
              <Button type="submit" className="w-full">
                Sign in with Spotify
              </Button>
            </Link>
          </CardFooter>
        </form>
      </Card>
    </section>
  );
}
