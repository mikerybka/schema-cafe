import React, { useState } from 'react';
import ReactDOM from 'react-dom/client';
import { Dialog } from '@headlessui/react'
import './dialog.css' // Make sure to import the CSS


const data = JSON.parse(document.getElementById('data')!.textContent!);
const path = window.location.pathname;

function filepathJoin(...parts: string[]): string {
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
    const [contents, setContents] = useState(props.contents);
    const [creating, setCreating] = useState(false);

    const refresh = () => {
        fetch(props.path, {
            headers: {
                'Accept': 'application/json',
            },
        }).then(res => res.json()).then(d => setContents(d.value));
    }

    return <>
        <ul>
            {contents.map(c => {
                return <li key={c.name}>
                    <a href={filepathJoin(props.path, id(c.name))}>{c.name}</a>
                </li>
            })}
        </ul>
        <button onClick={() => setCreating(true)}>New</button>
        <CreateModal isOpen={creating} onClose={() => {
            setCreating(false);
            refresh();
        }} />
    </>
}


function CreateModal({ isOpen, onClose }) {
    const [name, setName] = useState("")
    const [error, setError] = useState("")

    return (
        <Dialog open={isOpen} onClose={onClose} className="dialog-overlay">
            <div className="dialog-backdrop" aria-hidden="true" />
            <div className="dialog-container">
                <Dialog.Panel className="dialog-panel">
                    <Dialog.Title className="dialog-title"></Dialog.Title>
                    <Dialog.Description className="dialog-description">
                    </Dialog.Description>
                    <StringInput label='Name' value={name} onChange={setName} />
                    <div>{error}</div>
                    <button onClick={() => {
                        const p = filepathJoin(path, name)
                        fetch(p, {
                            method: "PUT",
                            body: JSON.stringify({ fields: [] }),
                        }).then(res => {
                            if (res.ok) {
                                window.location.pathname = p
                            } else {
                                res.text().then(t => setError(t))
                            }
                        }).catch(err => setError(err));
                    }}>Create</button>
                </Dialog.Panel>
            </div>
        </Dialog>
    )
}

function TitleBar() {
    return <div>{path}</div>
}

function Schema(props: {
    path: string;
    fields: {
        name: string;
        type: string;
    }[];
}) {
    const [error, setError] = useState("")
    const [fields, setFields] = useState(props.fields);
    const [saving, setSaving] = useState(false);
    const createField = () => {
        setFields([...fields, { name: "", type: "" }])
    }
    const deleteField = (index: number) => {
        setFields(fields => fields.filter((f, i) => i !== index));
    }
    const setFieldName = (index: number, name: string) => {
        setFields(fields => fields.map((f, i) => i === index ? { ...f, name } : f));
    }
    const setFieldType = (index: number, type: string) => {
        setFields(fields => fields.map((f, i) => i === index ? { ...f, type } : f));
    }
    const save = () => {
        setSaving(true);
        fetch(path, {
            method: "PUT",
            body: JSON.stringify({
                fields,
            })
        }).then(res => {
            if (res.ok) {
                setSaving(false);
            } else {
                console.log(res);
                res.text().then(text => {
                    setError(text);
                    console.log(text);
                });
            }
        })
    }
    return <div>
        <List title="Fields" onCreate={createField}>
            {fields.map((f, i) => {
                return <ListItem key={i} onDelete={() => deleteField(i)}>
                    <StringInput label='Name' value={f.name} onChange={name => setFieldName(i, name)} />
                    <StringInput label='Type' value={f.type} onChange={type => setFieldType(i, type)} />
                </ListItem>
            })}
        </List>
        <Button onClick={save} disabled={saving}>Save</Button>
        <div>{error}</div>
    </div>
}

function Button(props: {
    onClick: () => void;
    disabled?: boolean;
    children: any;
}) {
    return <button onClick={props.onClick} disabled={props.disabled}>{props.children}</button>
}

function List(props: { title: string; onCreate: () => void; children: any }) {
    return <div>
        <div>{props.title}</div>
        {props.children}
        <Button onClick={props.onCreate}>+</Button>
    </div>
}

function ListItem(props: {
    onDelete: () => void;
    children: any;
}) {
    return (
        <div style={{ position: 'relative', padding: '1rem', border: '1px solid #ccc' }}>
            <button
                onClick={props.onDelete}
                style={{
                    position: 'absolute',
                    top: '0.5rem',
                    right: '0.5rem',
                    background: 'transparent',
                    border: 'none',
                    fontSize: '1.25rem',
                    cursor: 'pointer',
                }}
                aria-label="Close"
            >
                Del
            </button>
            {props.children}
        </div>
    );
}


function StringInput(props: {
    label: string;
    value: string;
    onChange: (s: string) => void;
}) {
    return <div>
        <div>{props.label}:</div>
        <input type='text' value={props.value} onChange={e => props.onChange(e.target.value)} />
    </div>
}

function Data() {
    if (data.type === "dir") {
        return <Dir path={path} contents={data.value} />
    }
    if (data.type === "schema") {
        return <Schema path={path} fields={data.value.fields} />
    }
    return <div>{JSON.stringify(data)}</div>
}

const root = ReactDOM.createRoot(document.getElementById('root')!);
root.render(<div>
    <TitleBar />
    <Data />
</div>
);
