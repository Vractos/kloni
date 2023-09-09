import { Suspense } from "react";
import Table from "./table";
import TableSkeleton from "../components/skeleton/table_skeleton";
import Link from 'next/link';

export const runtime = "edge";

export default async function AnnouncementPage({
  searchParams,
}: {
  searchParams?: { [key: string]: string | string[] | undefined };
}) {
  const { q: searchValue } = searchParams as { [key: string]: string };

  return (
    <>
      <main>
        <div className="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
          <div className="px-4 py-6 sm:px-0 ">
            {!searchValue ? (
              <div className="px-4 py-6 sm:px-0">
                <p className="mt-5 text-center">
                  Busque por anúncios através do SKU
                </p>
              </div>
            ) : (
              <Suspense key={searchValue} fallback={<TableSkeleton />}>
                <Table sku={searchValue} />
              </Suspense>
            )}
          </div>
        </div>
      </main>
    </>
  );
}
