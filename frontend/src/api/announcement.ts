import announcements from '../mocks/announcements.json'
import { IAnnouncement } from "../interfaces/announcement";
import { api } from "./config";
import { errorMessages } from '../constants/messages';

export function getAnnouncements(sku: string): Promise<IAnnouncement[]> {
  const retrievedAnnouncements = announcements.filter(announcements => announcements.sku.toLowerCase() === sku)

  return new Promise((res, rej) => {
    setTimeout(() => {
      if(!retrievedAnnouncements.length) {
        rej({
          message: errorMessages.ANNOUNCEMENT_NOT_FOUND
        });
      } else {
        res(retrievedAnnouncements);
      }
    }, 1000);
    // api
    //   .get(`/announcements/${sku}`)
    //   .then((resp) => 
    //     res(
    //       resp.data.map((announcement: IAnnouncement): IAnnouncement => ({
    //         picture: announcement.picture,
    //         title: announcement.title,
    //         price: announcement.price,
    //         quantity: announcement.quantity,
    //         sku: announcement.sku
    //       }))
    //     )
    //   )
    //   .catch((err) => {
    //     if (err.status === 404) {
    //       rej({
    //         ...err,
    //         message: errorMessages.ANNOUNCEMENT_NOT_FOUND,
    //       });
    //     } else {
    //       rej(err);
    //     }
    //   });
  })
}