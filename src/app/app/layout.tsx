export const dynamic = 'force-dynamic';
import { ThemeProvider } from "@/components/ThemeProvider";
import {
  Home,
  LineChart,
  Menu,
  Package,
  Package2,
  Search,
  ShoppingCart,
  Users,
} from "lucide-react";
import Link from "next/link";

import { ModeToggle } from "@/components/ModeToggle";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { UserDropdown } from "@/components/UserDropdown";
import { cookies } from "next/headers";
import PocketBase from "pocketbase";
import { PlaylistRecord, UserRecord } from "@/interfaces";
import { SearchBar } from "@/components/SearchBar";

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const pb = new PocketBase(process.env.POCKETBASE_URL);
  await pb.admins.authWithPassword(
    process.env["ADMIN_USER"] ?? "",
    process.env["ADMIN_PASSWORD"] ?? ""
  );

  const cookie = cookies()?.get("pb_auth");
  const cookieData = JSON.parse(decodeURIComponent(cookie?.value ?? ""));
  const userId = cookieData.user_id;

  const userData: UserRecord = await pb
    .collection("users")
    .getFirstListItem(`username="${userId}"`);

  const playlistsResponse = await fetch(
    `${process.env["BACKEND_URL"]}/user/playlists/${userData.id}`
  );
  const playlists: PlaylistRecord[] = await playlistsResponse.json();
  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <div className="grid min-h-screen w-full md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr]">
        <div className="hidden border-r bg-muted/40 md:block">
          <div className="flex h-full max-h-screen flex-col gap-2">
            <div className="flex h-14 items-center border-b px-4 lg:h-[60px] lg:px-6">
              <Link href="/" className="flex items-center gap-2 font-semibold">
                <Package2 className="h-6 w-6" />
                <span className="">daylist rewind</span>
              </Link>
            </div>
            <div className="flex-1">
              <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
                <Link
                  href="#"
                  className="flex items-center gap-3 rounded-lg bg-muted px-3 py-2 text-primary transition-all hover:text-primary"
                >
                  <Package className="h-4 w-4" />
                  WIP
                </Link>
              </nav>
            </div>
            <div className="mt-auto p-4">
              <Card x-chunk="dashboard-02-chunk-0">
                <CardHeader className="p-2 pt-0 md:p-4">
                  <CardTitle>Donate</CardTitle>
                  <CardDescription>
                    Help support my work by donating.
                  </CardDescription>
                </CardHeader>
                <CardContent className="p-2 pt-0 md:p-4 md:pt-0">
                  <Link href={"http://jason-wong.me/donate"}>
                    <Button size="sm" className="w-full">
                      Donate
                    </Button>
                  </Link>
                </CardContent>
              </Card>
              created by jasonw
            </div>
          </div>
        </div>
        <div className="flex flex-col">
          <header className="flex h-14 items-center gap-4 border-b bg-muted/40 px-4 lg:h-[60px] lg:px-6">
            <Sheet>
              <SheetTrigger asChild>
                <Button
                  variant="outline"
                  size="icon"
                  className="shrink-0 md:hidden"
                >
                  <Menu className="h-5 w-5" />
                  <span className="sr-only">Toggle navigation menu</span>
                </Button>
              </SheetTrigger>
              <SheetContent side="left" className="flex flex-col">
                IT'S NOT MOBILE RESPONSIVE YET DONT LOOK
                {/* <nav className="grid gap-2 text-lg font-medium">
                  <Link
                    href="#"
                    className="flex items-center gap-2 text-lg font-semibold"
                  >
                    <Package2 className="h-6 w-6" />
                    <span className="sr-only">Acme Inc</span>
                  </Link>
                  <Link
                    href="#"
                    className="mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 text-muted-foreground hover:text-foreground"
                  >
                    <Home className="h-5 w-5" />
                    Dashboard
                  </Link>
                  <Link
                    href="#"
                    className="mx-[-0.65rem] flex items-center gap-4 rounded-xl bg-muted px-3 py-2 text-foreground hover:text-foreground"
                  >
                    <ShoppingCart className="h-5 w-5" />
                    Orders
                    <Badge className="ml-auto flex h-6 w-6 shrink-0 items-center justify-center rounded-full">
                      6
                    </Badge>
                  </Link>
                  <Link
                    href="#"
                    className="mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 text-muted-foreground hover:text-foreground"
                  >
                    <Package className="h-5 w-5" />
                    Products
                  </Link>
                  <Link
                    href="#"
                    className="mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 text-muted-foreground hover:text-foreground"
                  >
                    <Users className="h-5 w-5" />
                    Customers
                  </Link>
                  <Link
                    href="#"
                    className="mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 text-muted-foreground hover:text-foreground"
                  >
                    <LineChart className="h-5 w-5" />
                    Analytics
                  </Link>
                </nav> */}
                <div className="mt-auto">
                  <Card>
                    <CardHeader>
                      <CardTitle>BUT YOU CAN STILL DONATE</CardTitle>
                      <CardDescription>
                        CLICK THIS BUTTON TO SEND ME MONEY
                      </CardDescription>
                    </CardHeader>
                    <CardContent>
                      <Link href="https://jason-wong.me/donate">
                        <Button size="sm" className="w-full">
                          DONATE
                        </Button>
                      </Link>
                    </CardContent>
                  </Card>
                </div>
              </SheetContent>
            </Sheet>
            <div className="w-full flex-1">
              <form>
                <div className="relative">
                  <SearchBar playlists={playlists} />
                  {/* <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                  <Input
                    type="search"
                    placeholder="Search playlists..."
                    className="w-full appearance-none bg-background pl-8 shadow-none md:w-2/3 lg:w-1/3"
                  /> */}
                </div>
              </form>
            </div>
            {userData?.display_name}
            <UserDropdown avatar_url={userData?.avatar_url ?? ""} />
            <ModeToggle />
          </header>
          <main className="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
            {children}
            <div className="flex flex-col items-end w-full text-muted-foreground">
              <span>created by jasonw</span>{" "}
              <span>not affiliated with spotify</span>
            </div>
          </main>
        </div>
      </div>
    </ThemeProvider>
  );
}
