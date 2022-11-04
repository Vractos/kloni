import React, { FormEvent, Fragment, useRef, useState } from 'react'
import { Dialog, Transition } from '@headlessui/react'
import { XMarkIcon } from '@heroicons/react/24/outline'

interface ModalProps {
  isOpen: boolean;
  handleClose: React.Dispatch<React.SetStateAction<boolean>>

}

interface IInputs {
  // TODO: Add key type (string)
  title: string
}
const Modal: React.FC<ModalProps> = ({ isOpen, handleClose }) => {
  // const cancelButtonRef = useRef(null)
  const [inputFields, setInputField] = useState<IInputs[]>([{ title: '' }])

  const cancelButtonRef = useRef(null)

  function addInput(e: FormEvent) {
    e.preventDefault();
    let newField = { title: '' }

    setInputField([...inputFields, newField])
  }

  function removeInput(e: FormEvent, index: number) {
    e.preventDefault();
    let data = [...inputFields]
    data.splice(index, 1)
    setInputField(data)
  }

  function closeModal(): void {
    handleClose(false)
    setTimeout(() => {
      setInputField([{ title: '' }])
    }, 650);
  }

  function handleFormChange(index: number, event: React.ChangeEvent<HTMLInputElement>) {
    let data = [...inputFields];
    data[index].title = event.target.value
    setInputField(data)
  }

  return (
    <Transition.Root appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-10" initialFocus={cancelButtonRef} onClose={() => closeModal()}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-black bg-opacity-25" />
        </Transition.Child>

        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex min-h-full items-center justify-center p-4 text-center">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 scale-95"
              enterTo="opacity-100 scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 scale-100"
              leaveTo="opacity-0 scale-95"
            >
              <Dialog.Panel className="w-full max-w-md transform overflow-hidden rounded-2xl bg-white p-6 text-left align-middle shadow-xl transition-all">
                <Dialog.Title
                  as="h3"
                  className="text-lg font-medium leading-6 text-gray-900"
                >
                  Clones
                </Dialog.Title>
                <div className="">
                  <form action="#" className='mt-5 sm:flex-col sm:items-center' method='POST'>
                    {inputFields.map((input, index) => {
                      return (
                        <div className='flex mt-3'>
                          <div className="w-full sm:max-w-xs">
                            <label htmlFor="title" className="sr-only">
                              Email
                            </label>
                            <input
                              type="text"
                              name="title"
                              id="title"
                              value={input.title}
                              onChange={e => handleFormChange(index, e)}
                              className="block w-full rounded-md min-h-full border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-2"
                              placeholder="Insira o título do anúncio"
                            />
                          </div>
                          {inputFields.length > 1 ?
                            <button
                              onClick={e => removeInput(e, index)}
                              className="mt-3 inline-flex w-full items-center justify-center rounded-md border border-transparent bg-red-600 px-4 py-2 font-medium text-white shadow-sm hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                            >
                            <XMarkIcon className='h-5 w-4 text-white'/>
                            </button>
                            :
                            <button
                              onClick={e => removeInput(e, index)}
                              className="mt-3 inline-flex w-full items-center justify-center rounded-md border border-transparent bg-red-600 px-4 py-2 font-medium text-white shadow-sm hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm invisible"
                            >
                            <XMarkIcon className='h-5 w-4 text-white'/>
                            </button>
                            
                          }
                        </div>
                      )
                    })}
                  </form>
                </div>

                <div className="mt-4 text-right">
                  <button
                    type="button"
                    className="inline-flex justify-center float-left rounded-md border border-transparent bg-blue-700 px-4 py-2 text-sm font-medium text-white hover:bg-blue-200 focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 mr-1"
                    onClick={e => addInput(e)}
                  >
                    +
                  </button>
                  <button
                    type="button"
                    className="inline-flex justify-center rounded-md border border-transparent bg-red-100 px-4 py-2 text-sm font-medium text-red-900 hover:bg-red-200 focus:outline-none focus-visible:ring-2 focus-visible:ring-red-500 focus-visible:ring-offset-2 mr-1"
                    onClick={() => closeModal()}
                    ref={cancelButtonRef}
                  >
                    Cancelar
                  </button>
                  <button
                    type="button"
                    className="inline-flex justify-center rounded-md border border-transparent bg-blue-100 px-4 py-2 text-sm font-medium text-blue-900 hover:bg-blue-200 focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2"
                  >
                    Clonar
                  </button>
                </div>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div >
      </Dialog >
    </Transition.Root >
  )
}

export default Modal;