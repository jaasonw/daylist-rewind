"use client";

import * as React from "react";
import { Music } from "lucide-react";

import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
  CommandShortcut,
} from "@/components/ui/command";
import { useRouter } from "next/navigation";
import { PlaylistRecord } from "@/interfaces";
import Link from "next/link";

export function SearchBar({ playlists }: { playlists: PlaylistRecord[] }) {
  const [open, setOpen] = React.useState(false);
  const router = useRouter();

  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "j" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpen((open) => !open);
      }
    };

    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  const toggle = () => setOpen((open) => !open);

  const runCommand = React.useCallback((command: () => unknown) => {
    setOpen(false);
    command();
  }, []);

  return (
    <>
      <button
        className="flex items-center whitespace-nowrap transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 border border-input hover:bg-accent hover:text-accent-foreground relative px-2 pl-4 h-8 w-full justify-between rounded-[0.5rem] bg-muted/50 text-sm font-normal text-muted-foreground shadow-none"
        onClick={toggle}
      >
        <span>Search Playlists</span>
        <p className="text-sm text-muted-foreground">
          Press{" "}
          <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
            <span className="text-xs">⌘</span>J
          </kbd>
        </p>
      </button>

      <CommandDialog open={open} onOpenChange={setOpen}>
        <CommandInput placeholder="Type a command or search..." />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          <CommandGroup heading="Latest">
            {playlists.map((playlist: PlaylistRecord) => (
              <CommandItem
                key={playlist.hash}
                value={playlist.title}
                onSelect={() => {
                  runCommand(() => router.push(`/app/playlist/${playlist.id}`));
                }}
              >
                <Music className="mr-2 h-4 w-4" />
                <span>{playlist.title.replaceAll("daylist • ", "")}</span>
              </CommandItem>
            ))}
          </CommandGroup>
        </CommandList>
      </CommandDialog>
    </>
  );
}
