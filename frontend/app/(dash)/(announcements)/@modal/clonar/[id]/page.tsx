"use client";

import Modal from "@/components/modal";
import { Fragment, useEffect, useState } from "react";
import { clone } from "@/app/actions";
import {
  CheckCircleIcon,
  PlusCircleIcon,
  XMarkIcon,
} from "@heroicons/react/24/outline";
import { useRouter } from "next/navigation";
import { Dialog, Transition } from "@headlessui/react";
import { experimental_useFormState as useFormState } from "react-dom";
import SubmitButton from "@/components/submitButton";
import Notification from "@/components/notification";
import { XCircleIcon } from '@heroicons/react/20/solid';

const initialState = {
  fails: 0,
};

export default function CloneForm({
  params,
  searchParams,
}: {
  params: { id: string };
  searchParams?: { [key: string]: string | string[] | undefined };
}) {
  const { q: sku } = searchParams as { [key: string]: string };

  const [inputs, setInputs] = useState([""]);
  const [open, setOpen] = useState(true);
  const [failMessage, setFailMessage] = useState(false);
  const [state, formAction] = useFormState(clone, initialState);
  const router = useRouter();

  useEffect(() => {
    if (state.fails) {
      setFailMessage(true);
      setTimeout(() => {
        setFailMessage(false);
      }, 2500);
    }
  }
  , [state.fails]);
  
  const addInput = () => {
    setInputs([...inputs, ""]);
  };

  const removeInput = (index: number) => {
    console.log(index);
    console.log(inputs);
    console.log(inputs[index]);
    if (inputs.length > 1) {
      let newInputs = [...inputs];
      newInputs.splice(index, 1);
      console.log(newInputs);
      setInputs(newInputs);
    }
  };

  const handleInputChange = (index: number, value: string) => {
    let newInputs = [...inputs];
    newInputs[index] = value;
    setInputs(newInputs);
  };

  const onClose = () => {
    setOpen(false);
    setTimeout(() => {
      router.back();
    }, 600);
  };

  return (
    <>
      <Notification showing={failMessage}>
        <div className="rounded-md bg-red-50 p-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <XCircleIcon
                className="h-5 w-5 text-red-400"
                aria-hidden="true"
              />
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-red-800">
                Falha ao clonar. Por favor, tente novamente.
              </h3>
            </div>
          </div>
        </div>
      </Notification>
      <Modal isOpen={open}>
        <form
          action={formAction}
          className="bg-white shadow-sm ring-1 ring-gray-900/5 sm:rounded-xl md:col-span-2 pt-6"
        >
          <input type="hidden" name="id" value={params.id} />
          <input type="hidden" name="sku" value={sku} />
          <div className="absolute right-0 top-0 hidden pr-4 pt-4 sm:block">
            <button
              type="button"
              className="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
              onClick={onClose}
            >
              <span className="sr-only">Close</span>
              <XMarkIcon className="h-6 w-6" aria-hidden="true" />
            </button>
          </div>
          <div className="mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left  px-4">
            <Dialog.Title
              as="h2"
              className="text-base font-semibold leading-6 text-gray-900"
            >
              Clones
            </Dialog.Title>
            <div className="mt-2">
              <p className="text-sm text-gray-500">
                Adicione os títulos dos anúncios. Será criado um novo anúncio
                para cada título, com o mesmo SKU e preço do anúncio original.
              </p>
            </div>
          </div>
          {inputs.map((value, index) => (
            <div className="px-8 pt-6 sm:px-8" key={index}>
              <div className="flex gap-x-6">
                <label htmlFor="email-address" className="sr-only">
                  Título
                </label>
                <input
                  id="title"
                  name="title"
                  type="text"
                  maxLength={60}
                  required
                  value={value}
                  onChange={(e) => handleInputChange(index, e.target.value)}
                  className="min-w-0 flex-auto rounded-md border-0 px-3.5 py-2 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                  placeholder="Título do anúncio"
                />

                {inputs.length > 1 && (
                  <button
                    type="button"
                    onClick={() => removeInput(index)}
                    className="flex-none rounded-md bg-red-500 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-red-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-500 text-center"
                  >
                    <XMarkIcon className="h-5 w-5 text-white stroke-2" />
                  </button>
                )}
              </div>
            </div>
          ))}
          <div className="flex items-center justify-between border-t border-gray-900/10 px-4 py-4 sm:px-8 mt-6">
            <button
              type="button"
              className="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 self-start"
              onClick={addInput}
            >
              <PlusCircleIcon className="h-5 w-5 text-white stroke-1" />
            </button>
            <div className="flex items-center gap-x-6">
              <button
                type="button"
                className="text-sm font-semibold leading-6 text-gray-900"
                onClick={onClose}
              >
                Cancelar
              </button>
              <SubmitButton text="Clonar" />
            </div>
          </div>
        </form>
      </Modal>
    </>
  );
}
