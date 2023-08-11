"use client"; // This is a client component 👈🏽

import {
  createContext,
  memo,
  useState,
  useId,
} from 'react';
import Image from 'next/image'
import { IconSearch } from './iconsearch';

const fetcher = (...args) => fetch(...args).then(res => res.json())


interface SearchInputProps {
  value: string;
  // onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  onChange: (txt: string) => void;
  onSubmit: () => void;
}
function SearchInput({ value, onChange, onSubmit }: SearchInputProps) {
  const id = useId();
  return (
    <form
      className="mb-3 py-1"
      data-hover="SearchInput"
      onSubmit={(e) => e.preventDefault()}>
      <label htmlFor={id} className="sr-only">
        Search
      </label>
      <div className="relative w-full">
        <div className="absolute inset-y-0 left-0 flex items-center pl-4 pointer-events-none">
          <IconSearch className="text-gray-30 w-4" />
        </div>
        <input
          type="text"
          id={id}
          className="flex pl-11 py-4 h-10 w-full bg-secondary-button outline-none betterhover:hover:bg-opacity-80 pointer items-center text-left text-primary rounded-full align-middle text-base"
          placeholder="Search"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              onSubmit();
            }
          }} />
      </div>
    </form>
  );
}

function SyncButton() {

  const onClick = async () => {
    try {
      const res = await fetch('http://localhost:8080/sync')
      const data = await res.json();
      window.location.href = data.auth_url;
    }
    catch (e) {
      console.log('caught an error');
      console.log(e);
    }
  }

  return (
    <div className="flex flex-wrap -mx-3 mb-2 p-24 justify-between items-center">
      <div className="w-full md:w-1/2 px-3 mb-6 md:mb-0">
        <div className="absolute">
          <button onClick={() => onClick()} className="-transparent hover:bg-blue-500 text-blue-700 font-semibold hover:text-white py-2 px-4 border border-blue-500 hover:border-transparent rounded mx-8" type="button">
            Sync
          </button>
          <button onClick={() => onClick()} className="-transparent hover:bg-blue-500 text-blue-700 font-semibold hover:text-white py-2 px-4 border border-blue-500 hover:border-transparent rounded" type="button">
            Reload songs
          </button>
        </div>
      </div>
    </div>
  );
}

interface Song {
  title: string;
  content: string;
}
interface TableComponentProps {
  songs: Song[];
}
function TableComponent({ songs }: TableComponentProps) {
  console.log('rendering');
  console.log(songs);
  return (
    <table className="table-auto border-collapse border border-slate-400 border-spacing-2">
      <thead>
        <tr>
          <th className="px-4 py-2">Title</th>
          <th className="px-4 py-2">Singer</th>
          <th className="px-4 py-2">Lyrics</th>
        </tr>
      </thead>
      <tbody>
        {songs && songs.map((song, index) => {
          return (
            <tr key={index}>
              <td className="border px-4 py-2">{song.title}</td>
              <td className="border px-4 py-2">SINGER</td>
              <td className="border px-4 py-2">{song.content}</td>
            </tr>
          )
        })
        }
      </tbody>
    </table>
  );
}

export default function Home() {
  const [searchText, setSearchText] = useState('');
  const [songs, setSongs] = useState(new Array<Song>());
  const doSearch = async () => {
    console.log('searching for ' + searchText);

    const searchURL = new URL('http://localhost:8080/search');
    searchURL.searchParams.append('q', searchText);

    const res = await fetch(searchURL)

    const data:Song[] = await res.json();
    console.table(data);
    setSongs(data)
  }
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <SyncButton />
      <SearchInput
        value={searchText}
        onChange={newText => setSearchText(newText)}
        onSubmit={() => doSearch()}
      />

      <TableComponent songs={songs} />

    </main>
  )
}
