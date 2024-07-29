"use client";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Check, CircleX, Loader2, Plus } from "lucide-react";

export function ExportPlaylistButton() {
  const [exporting, setExporting] = useState(false);
  const [error, setError] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const [success, setSuccess] = useState(false);

  const handleClick = async () => {
    setExporting(true);
    setError(false);
    setSuccess(false);
    setErrorMessage("");

    try {
      const response = await fetch(window.location.href + "/export", {
        method: "GET",
      });
      if (response.ok) {
        console.log(response.text());
        setSuccess(true);
      } else {
        setError(true);
        setErrorMessage("Failed to export playlist.");
      }
    } catch (e) {
      setError(true);
      setErrorMessage("An error occurred.");
    } finally {
      setExporting(false);
    }
  };

  const icon = () => {
    if (exporting) {
      return (
        <>
          <Loader2 className="mr-2 h-4 w-4 animate-spin" size={24} />
          Exporting...
        </>
      );
    }
    if (error) {
      return (
        <>
          <CircleX className="mr-2" size={24} />
          {errorMessage}
        </>
      );
    }
    if (success) {
      return (
        <>
          <Check className="mr-2" size={24} />
          Success
        </>
      );
    }
    return (
      <>
        <Plus className="mr-2" size={24} />
        Save to Spotify
      </>
    );
  };

  return (
    <Button disabled={exporting} onClick={handleClick}>
      {icon()}
    </Button>
  );
}
