"use client";

import { MagnifyingGlassIcon } from "@heroicons/react/24/outline";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { createUrl } from "@/lib/utils/url";
import { Input } from "./ui/input";

export default function Search() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [searchValue, setSearchValue] = useState("");

  useEffect(() => {
    setSearchValue(searchParams?.get("q") || "");
    console.log(searchParams);
  }, [searchParams, setSearchValue]);

  function onSubmit(e: React.FormEvent) {
    e.preventDefault();

    const val = e.target as HTMLFormElement;
    const search = val.search as HTMLInputElement;
    const newParams = new URLSearchParams(searchParams.toString());

    if (search.value) {
      newParams.set("q", search.value);
    } else {
      newParams.delete("q");
    }

    router.push(createUrl("/anuncios", newParams));
  }

  return (
    <form onSubmit={onSubmit}>
      <Input
        type="search"
        name="search"
        value={searchValue}
        placeholder="Utilize o SKU para buscar..."
        onChange={(e) => setSearchValue(e.target.value)}
        onSubmit={(e) => onSubmit(e)}
        className="w-full rounded-lg bg-background pl-8 md:w-[200px] lg:w-[320px]"
      />
    </form>
  );
}
