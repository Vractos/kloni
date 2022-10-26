import React, { FormEvent, Fragment, useRef, useState } from 'react'
import { Dialog, Transition } from '@headlessui/react'
import { ExclamationTriangleIcon } from '@heroicons/react/24/outline'

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

  function addInput(e: FormEvent) {
    e.preventDefault();
    let newField = {title: ''}

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
    setInputField([{title: ''}])
  }

  function handleFormChange(index: number, event: React.ChangeEvent<HTMLInputElement>) {
    let data = [...inputFields];
    data[index].title = event.target.value
    setInputField(data)
  }

  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-10" onClose={handleClose}>
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
                <div className="mt-2">
                  <form action="#" method='POST'>
                    {inputFields.map((input, index) => {
                      return (
                        <div key={index} className="relative row-span-6 sm:row-span-3 mt-4" >
                          <input
                            type="text"
                            name="title"
                            id="title"
                            placeholder='Insira o título do anúncio'
                            value={input.title}
                            onChange={e => handleFormChange(index,e)}
                            autoComplete="family-name"
                            className=" outline-none block w-10/12 rounded-md border-gray-300 shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                          />
                          {inputFields.length > 1 && <button onClick={e => removeInput(e, index)} className="text-white absolute right-2.5 bottom-2 bg-red-600 hover:bg-red-700 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-2 py-1">x</button>}
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
    </Transition >
  )
}

export default Modal;