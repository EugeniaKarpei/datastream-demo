import React from "react"
import { FaGithub } from "react-icons/fa"

export default function Header(){
    return (
        <header>
            <a href="https://github.com/vkarpei/valery-datadog-datastream-demo" target="_blank" rel="noopener noreferrer">
                <FaGithub className="git_icon"/>
            </a>
        </header>
    )
}