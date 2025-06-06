// Función para navegar entre páginas
function navigateToPage(page) {
  console.log(`Navegando a: ${page}`);

  // Mapeo de páginas a archivos HTML
  const pageFiles = {
    multas: "multas.html",
    informacion: "informacion.html",
    configuracion: "configuracion.html",
  };

  // Verificar si la página existe
  if (!pageFiles[page]) {
    console.error(`Página no encontrada: ${page}`);
    return;
  }

  // Cargar el contenido de la página
  fetch(pageFiles[page])
    .then((response) => {
      if (!response.ok) {
        throw new Error(`Error al cargar la página: ${response.status}`);
      }
      return response.text();
    })
    .then((html) => {
      // Reemplazar todo el contenido del body
      document.body.innerHTML = html;

      // Re-ejecutar scripts si es necesario
      executePageScripts(page);
    })
    .catch((error) => {
      console.error("Error al navegar:", error);
      alert("Error al cargar la página. Por favor, intenta de nuevo.");
    });
}

// Función para ejecutar scripts específicos de cada página
function executePageScripts(page) {
  if (page === "multas") {
    setTimeout(() => {
      startMultasMonitoring();
    }, 100);
  }
}

// Sistema de multas INICIALIZAR
let multasData = [];
let multasInterval;

function startMultasMonitoring() {
  console.log("Iniciando monitoreo de multas...");

  const tableBody = document.getElementById("multas-table-body");
  const multasCount = document.getElementById("multas-count");

  if (!tableBody || !multasCount) {
    console.error("Elementos del DOM no encontrados. Reintentando en 500ms...");
    setTimeout(startMultasMonitoring, 100);
    return;
  }

  // Limpiar intervalo anterior si existe
  if (multasInterval) {
    clearInterval(multasInterval);
  }

  // Inicializar la tabla
  updateMultasTable();

  // Función para obtener datos y procesar multas
  function checkForMultas() {
    fetch("/_/info")
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        console.log("Datos recibidos:", data);

        document.querySelector("#probability").innerText =
          `${data.Probability}%`;

        //  puntos son mayor a 50, es una multa
        if (data.Probability >= 50) {
          console.log(`Multa Puntos: ${data.Probability}`);

          const nuevaMulta = {
            id: data.ID,
            fecha: new Date().toLocaleDateString("es-ES"),
            hora: new Date().toLocaleTimeString("es-ES", {
              hour: "2-digit",
              minute: "2-digit",
              second: "2-digit",
            }),
            infraccion: getInfractionType(data.Probability),
            puntos: data.Probability,
            monto: `$${getInfractionCost(data.Probability)}`,
          };

          let existeMultaSimilar = false;

          for (let multa of multasData) {
            if (multa.id === data.ID) {
              existeMultaSimilar = true;
              break;
            }
          }

          if (!existeMultaSimilar) {
            // Agregar la multa al inicio del array
            multasData.unshift(nuevaMulta);

            // Mantener solo las últimas 50 multas para mejor historial
            if (multasData.length > 50) {
              multasData = multasData.slice(0, 50);
            }

            console.log(`Multa agregada. Total multas: ${multasData.length}`);

            // Actualizar la tabla
            updateMultasTable();
          } else {
            console.log("Multa similar ya existe, no se agrega duplicado");
          }
        }
      })
      .catch((err) => {
        console.error("Error al obtener datos:", err);
      });
  }

  // Ejecutar inmediatamente y luego cada 3 segundos (reducir frecuencia para evitar duplicados)
  checkForMultas();
  multasInterval = setInterval(checkForMultas, 500);
}

function getInfractionType(points) {
  if (points >= 90) return "Exceso de velocidad Demasiado Grave";
  if (points >= 80) return "Exceso de velocidad Grave";
  if (points >= 70) return "Exceso de velocidad Moderado";
  if (points >= 60) return "Exceso de velocidad Bajo";
  return "Exceso de velocidad Bajo";
}

function getInfractionCost(probability) {
  if (probability >= 90) return 1200000;
  if (probability >= 80) return 900000;
  if (probability >= 70) return 500000;
  if (probability >= 60) return 450000;
  return 300000;
}

function updateMultasTable() {
  const tableBody = document.getElementById("multas-table-body");
  const multasCount = document.getElementById("multas-count");

  if (!tableBody) {
    console.error("Elemento 'multas-table-body' no encontrado");
    return;
  }

  // Limpiar tabla
  tableBody.innerHTML = "";

  if (multasData.length === 0) {
    tableBody.innerHTML = `
      <tr>
        <td colspan="6" class="px-4 py-8 text-center text-gray-500">
          <i class="fa-solid fa-clock mr-2"></i>
          Esperando datos de multas...
        </td>
      </tr>
    `;
  } else {
    // Agregar cada multa a la tabla
    multasData.forEach((multa, index) => {
      const row = document.createElement("tr");
      row.className =
        index % 2 === 0 ? "bg-white border-b" : "bg-gray-50 border-b";
      row.innerHTML = `
        <td class="px-4 py-3 text-sm">${multa.fecha}</td>
        <td class="px-4 py-3 text-sm">${multa.hora}</td>
        <td class="px-4 py-3 text-sm">${multa.infraccion}</td>
        <td class="px-4 py-3 text-sm">
          <span class="bg-red-100 text-red-800 px-2 py-1 rounded-full text-xs font-medium">
            ${multa.puntos} %
          </span>
        </td>
        <td class="px-4 py-3 text-sm font-semibold text-green-600">${multa.monto}</td>
      `;
      tableBody.appendChild(row);
    });
  }

  // Actualizar contador de multas
  if (multasCount) {
    multasCount.textContent = multasData.length;
    console.log(`Contador actualizado: ${multasData.length} multas`);
  } else {
    console.error("Elemento 'multas-count' no encontrado");
  }
}

// Función para volver al inicio
function goHome() {
  if (multasInterval) {
    clearInterval(multasInterval);
    multasInterval = null;
    console.log("Monitoreo de multas detenido");
  }

  window.location.href = "/";
}
