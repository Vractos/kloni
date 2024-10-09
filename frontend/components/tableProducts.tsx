import Image from "next/image"
import { MoreHorizontal } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { getAnnouncements } from '../app/api/handlers/announcements'
import { Status } from '../lib/interfaces/announcements'

export default async function AnnouncementsTable({
  sku,
}: {
  sku: string;
}) {
  const announcements = await getAnnouncements(sku)

  const tableHeaderClasses = "sticky top-0 bg-white z-10"
  const columnClasses = {
    image: "w-[100px] hidden sm:table-cell",
    title: "w-[250px]",
    sku: "w-[120px]",
    status: "w-[100px]",
    store: "w-[120px] hidden md:table-cell",
    quantity: "w-[100px] hidden md:table-cell",
    price: "w-[100px]",
    actions: "w-[50px]",
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Anúncios</CardTitle>
        <CardDescription>Gerencie seus anúncios e seus clones.</CardDescription>
      </CardHeader>
      <CardContent className="p-0">
        <div className="overflow-hidden rounded-md border">
          <div className="max-h-[73vh] overflow-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className={`${tableHeaderClasses} ${columnClasses.image}`}>
                    <span className="sr-only">Imagem</span>
                  </TableHead>
                  <TableHead className={`${tableHeaderClasses} ${columnClasses.title}`}>Título</TableHead>
                  <TableHead className={`${tableHeaderClasses} ${columnClasses.sku}`}>Sku</TableHead>
                  <TableHead className={`${tableHeaderClasses} ${columnClasses.status}`}>Status</TableHead>
                  <TableHead className={`${tableHeaderClasses} ${columnClasses.store}`}>Loja</TableHead>
                  <TableHead className={`${tableHeaderClasses} ${columnClasses.quantity}`}>Quantidade</TableHead>
                  <TableHead className={`${tableHeaderClasses} ${columnClasses.price}`}>Preço</TableHead>
                  <TableHead className={`${tableHeaderClasses} ${columnClasses.actions}`}>
                    <span className="sr-only">Ações</span>
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {announcements && announcements.map((announcement, index) => (
                  <TableRow key={index}>
                    <TableCell className={columnClasses.image}>
                      <Image
                        alt={announcement.title}
                        className="aspect-square rounded-md object-cover"
                        height="64"
                        src={announcement.picture}
                        width="64"
                      />
                    </TableCell>
                    <TableCell className={`font-medium ${columnClasses.title}`}>
                      <a href={announcement.link} target="_blank" rel="noreferrer">
                        {announcement.title}
                      </a>
                    </TableCell>
                    <TableCell className={columnClasses.sku}>{announcement.sku}</TableCell>
                    <TableCell className={columnClasses.status}>
                      <Badge 
                        className={announcement.status === Status.Ativo ? 'bg-green-500' : ''} 
                        variant={announcement.status === Status.Pausado ? "secondary" : "default"}
                      >
                        {announcement.status}
                      </Badge>
                    </TableCell>
                    <TableCell className={columnClasses.store}>
                      {announcement.account.name}
                    </TableCell>
                    <TableCell className={columnClasses.quantity}>
                      {announcement.quantity}
                    </TableCell>
                    <TableCell className={columnClasses.price}>
                      {announcement.price}
                    </TableCell>
                    <TableCell className={columnClasses.actions}>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button aria-haspopup="true" size="icon" variant="ghost">
                            <MoreHorizontal className="h-4 w-4" />
                            <span className="sr-only">Toggle menu</span>
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuLabel>Actions</DropdownMenuLabel>
                          <DropdownMenuItem>Edit</DropdownMenuItem>
                          <DropdownMenuItem>Delete</DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </div>
        {announcements && announcements.length === 0 && (
          <div className="px-4 py-5 sm:p-6 text-center">
            <h3 className="mt-3 text-base font-semibold text-gray-900">
              Nada por aqui...
            </h3>
            <p className="mt-3 text-xs text-gray-500">
              Você não possui anúncios com esse SKU.
            </p>
          </div>
        )}
      </CardContent>
      <CardFooter>
        <div className="text-xs pt-2 text-muted-foreground">
          Mostrando <strong>1-{announcements.length}</strong> de <strong>{announcements.length}</strong> anúncios
        </div>
      </CardFooter>
    </Card>
  )
}