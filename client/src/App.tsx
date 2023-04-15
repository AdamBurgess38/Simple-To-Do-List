import { useState } from 'react'
import reactLogo from './assets/react.svg'
import './App.css'
import { Box, List, ThemeIcon } from "@mantine/core"
import React from 'react';
import useSWR from "swr";
import AddToDo from './components/AddToDo';
import { CheckCircleFillIcon } from '@primer/octicons-react';
import { ListItem } from '@mantine/core/lib/List/ListItem/ListItem';

export interface ToDo{
  ID: number;
  Title: string;
  Body: string;
  Done: boolean;
} 

export const ENDPOINTGOLANG = "http://localhost:8080";

const fetcher = (url: string) =>
  fetch(`${ENDPOINTGOLANG}/${url}`).then((r) => r.json());

function App() {
  let headers = new Headers();
  headers.append('Content-Type', 'application/json');
  headers.append('Accept', 'application/json');

  // headers.append('Access-Control-Allow-Credentials', 'true');



  fetch(ENDPOINTGOLANG+'/healthcheck',{
    headers: headers})
    .then(response => {
      // console.log(response.json())
      const reader = response.body?.getReader();
      const decoder = new TextDecoder();
      if(!reader){ return}
      return reader.read().then(function processText({ done, value }) : any{
        if (done) {
          return;
        }
  
        const data = decoder.decode(value);
        console.log(data);
  
        return JSON.parse(data);}
      )})
    .then(data => {
      console.log(data)
    })
    .catch(error => {
      console.error('Error:', error);
    });

  // const { data, mutate } = useSWR<ToDo[]>("api/", fetcher);
  
  // async function markTodoAdDone(id: number) {
  //   const updated = await fetch(`${ENDPOINTGOLANG}/api/todos/${id}/done`, {
  //     method: "PATCH",
  //   }).then((r) => r.json());

  //   console.log(updated)
  //   mutate(updated);
  // }

  // return (
  //   <Box
  //     sx={(theme) => ({
  //       padding: "2rem",
  //       width: "100%",
  //       maxWidth: "40rem",
  //       margin: "0 auto",
  //     })}
  //   >
  //     <List spacing="xs" size="sm" mb={12} center>
  //       {data?.map((todo) => {
  //         console.log(todo)
  //         console.log(todo.ID)
  //         console.log(todo.Title)
  //         console.log(todo.Body)
  //         return (
  //           <List.Item
  //             key={`todo_list__${todo.ID}`}
  //             onClick={() => markTodoAdDone(todo.ID)}
  //             icon={
  //               todo.Done ? (
  //                 <ThemeIcon color="teal" size={24} radius="xl">
  //                   <CheckCircleFillIcon size={20} />
  //                 </ThemeIcon>
  //               ) : (
  //                 <ThemeIcon color="gray" size={24} radius="xl">
  //                   <CheckCircleFillIcon size={20} />
  //                 </ThemeIcon>
  //               )
  //             }
  //           >
  //             {todo.Body}
  //           </List.Item>
  //         );
  //       })}
  //     </List>

  //     <AddToDo mutate={mutate} />
  //   </Box>
  // );
}

export default App;