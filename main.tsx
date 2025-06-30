import React, { useState } from 'react';
import ReactDOM from 'react-dom/client';

const data = JSON.parse(document.getElementById('data')!.textContent!);
const path = window.location.pathname;

function joinPath(...parts: string[]): string {
    return parts
        .filter(Boolean)
        .map((part, index) => {
            // Remove leading slashes on everything but the first
            if (index > 0) part = part.replace(/^\/+/, '');
            // Remove trailing slashes on everything but the last
            if (index < parts.length - 1) part = part.replace(/\/+$/, '');
            return part;
        })
        .join('/');
}

function id(s: string) {
    return s
}


function Dir(props: { path: string; contents: { name: string; type: string }[] }) {
    console.log(path)
    return <ul>
        {props.contents.map(c => {
            return <li key={c.name}>
                <a href={joinPath(props.path + id(c.name))}>{c.name}</a>
            </li>
        })}
    </ul>
}

function TitleBar(props:{path: string}) {
    return <div>{props.path}</div>
}

function Schema(props: {
    path: string;
    fields: {
        name: string;
        type: string;
    }[];
}) {
    const [fields, setFields] = useState(props.fields);
    const createField = () => {
        setFields([...fields, {name:"", type:""}])
    }
    const deleteField = (index: number) => {
        setFields(fields => fields.filter((f, i) => i !== index));
    }
    return <div>
        <TitleBar path={props.path} />
        <List title="Fields" onCreate={createField}>
            {fields.map((f, i) => {
                return <ListItem key={i} onDelete={() => deleteField(i)}>
                    <StringInput label='Name' value={f.name} onChange={name => setFieldName(i, name)} />
                    <StringInput label='Type' value={f.type} onChange={type => setFieldType(i, type)} />
                </ListItem>
            })}
        </List>
    </div>
}

function Field(props: {
    name: string;
    setName: (name: string)=> void;
    type: string;
    setType: (type: string)=> void;
}) {

}

function StringInput(props: {
    label: string;
    value: string;
    onChange: (s: string) => void;
}) {

}

function App() {
    if (data.type === "dir") {
        return <Dir path={path} contents={data.value} />
    }
    if (data.type === "schema") {
        return <Schema path={path} fields={data.value.fields} />
    }
    return <div>{JSON.stringify(data)}</div>
}

const root = ReactDOM.createRoot(document.getElementById('root')!);
root.render(<App />);
