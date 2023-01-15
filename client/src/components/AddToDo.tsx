import { useState } from 'react'
import { useForm } from "@mantine/form";
import { Button, Group, Modal, Textarea, TextInput } from '@mantine/core'
import React from 'react'
import { ENDPOINTGOLANG, ToDo } from '../App';
import { KeyedMutator } from 'swr';

function AddToDo({mutate}: { mutate: KeyedMutator<ToDo[]> }){
    const [open, setOpen] = useState(false)


    const form = useForm({
        initialValues:{
            title: "",
            body: "",
        },
    })

    async function createToDo(values: {title: string, body: string}){
        const updated = await fetch(`${ENDPOINTGOLANG}/api/todos`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(values),
          },).then((r) => r.json());
      
          mutate(updated)
          form.reset()
          setOpen(false)
    
    }


    return <>
        <Modal opened={open} onClose={() => setOpen(false)} title="Create ToDo">
            <form onSubmit={form.onSubmit(createToDo)}>
                <TextInput 
                required mb={12} 
                label="ToDo"
                placeholder='What do you want to do?'
                {...form.getInputProps("title")}/>
                <Textarea 
                required mb={12} 
                label="Body"
                placeholder='Tell Me More...'
                {...form.getInputProps("body")}/>
                <Button type="submit">Create To Do</Button>
            </form>
        </Modal>

        <Group position="center">
            <Button fullWidth mb={12} onClick={() => setOpen(true)}
            ></Button>
        </Group>
    </>

}


export default AddToDo

