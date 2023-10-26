import SideMenu from './setting-nav';

export default function SettingLayout({ children }: { children: React.ReactNode}) { 

  return (
    <div className="max-h-screen" style={{ overflow: "auto" }}>
      <div className='max-w-7xl mx-auto lg:flex lg:gap-x-16 lg:px-8'>

      <SideMenu />
      {children}
      </div>
    </div>
  )
}