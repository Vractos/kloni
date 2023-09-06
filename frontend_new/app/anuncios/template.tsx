import SearchBar from '../components/searchBar';

export default function AnnouncementPage({children}: {children: React.ReactNode}) {
  return (
      <div className="mx-auto max-w-7xl py-6 px-4 sm:px-6 lg:px-8 ">
        <SearchBar />
        {children}
      </div>
  );
}
