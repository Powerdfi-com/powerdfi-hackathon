"use client"
import React, { ReactNode, useRef } from 'react'

const FileUpload = ({ children, updateValue, multiple = false }: { children: ReactNode, updateValue: (file: FileList) => void, multiple?: boolean }) => {
    const ref = useRef<HTMLInputElement>(null);
    return (
        <div onClick={() => ref.current?.click()} className="cursor-pointer z-10">
            <input type='file' ref={ref} hidden multiple={multiple} onChange={(e) => updateValue(e.target.files!)} />
            {children}
        </div>
    )
}

export default FileUpload